package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dsmit05/metida/internal/logger"
	"github.com/Dsmit05/metida/internal/models"
	"github.com/Dsmit05/metida/internal/repositories/postgres"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
)

type DBConnectI interface {
	GetConnectDB() string
}

// PostgresRepository repository for accessing Postgres database
type PostgresRepository struct {
	conn    *pgx.Conn
	queries *postgres.Queries
}

func NewPostgresRepository(url DBConnectI) (*PostgresRepository, error) {
	conn, err := pgx.Connect(context.Background(), url.GetConnectDB())
	if err != nil {
		return nil, err
	}
	queries := postgres.New(conn)

	return &PostgresRepository{
		conn:    conn,
		queries: queries,
	}, nil
}

func (o *PostgresRepository) Close() {
	err := o.conn.Close(context.Background())
	if err != nil {
		logger.L.Error("(o *PostgresRepository) Close() error:", err)
	}
}

func (o *PostgresRepository) CreateUser(name, password, email, role string) error {
	err := o.queries.CreateUser(context.Background(), postgres.CreateUserParams{
		Name:     sql.NullString{name, true},
		Password: password,
		Email:    email,
		Role:     role,
	})

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		err = fmt.Errorf("Пользователь уже существует")
		return err
	}

	if err != nil {
		return fmt.Errorf("что то пошло не так")
	}

	return nil
}

func (o *PostgresRepository) ReadUser(email string) (*models.User, error) {
	user, err := o.queries.ReadUser(context.Background(), email)

	if err != nil {
		if user.Email == "" {
			return nil, fmt.Errorf("Пользователя не существует")
		}

		return nil, fmt.Errorf("что то пошло не так")
	}

	if user.IsDeleted.Bool {
		return nil, fmt.Errorf("Пользователя не существует, или аккаунт был удален")
	}

	userModel := &models.User{
		ID:        user.ID,
		Name:      user.Name.String,
		Password:  user.Password,
		Email:     user.Email,
		Role:      user.Role,
		IsDeleted: user.IsDeleted.Bool,
	}

	return userModel, nil
}

func (o *PostgresRepository) UpdateUser(email string, name, password, role string, isDeleted bool) error {
	var err error
	err = o.queries.UpdateUser(context.Background(), postgres.UpdateUserParams{
		Email:     email,
		Name:      sql.NullString{name, true},
		Password:  password,
		Role:      role,
		IsDeleted: sql.NullBool{isDeleted, true},
	})

	return err
}

func (o *PostgresRepository) DeleteUser(email string) error {
	var err error
	err = o.queries.DeleteUser(context.Background(), email)

	return err
}

func (o *PostgresRepository) CreateSession(email string, refreshToken, userAgent, ip string, expiresIn int64) error {
	err := o.queries.CreateSession(context.Background(), postgres.CreateSessionParams{
		UserEmail:    sql.NullString{email, true},
		RefreshToken: sql.NullString{refreshToken, true},
		AccessToken:  sql.NullString{Valid: false},
		UserAgent:    sql.NullString{userAgent, true},
		Ip:           sql.NullString{ip, true},
		ExpiresIn:    expiresIn,
	})

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		err = fmt.Errorf("Сессия уже существует")
		return err
	}

	if err != nil {
		return fmt.Errorf("что то пошло не так")
	}

	return nil
}

func (o *PostgresRepository) ReadSession(email string, userAgent, ip string) (*models.Session, error) {
	var err error

	session, err := o.queries.ReadSession(context.Background(), postgres.ReadSessionParams{
		UserEmail: sql.NullString{email, true},
		UserAgent: sql.NullString{userAgent, true},
		Ip:        sql.NullString{ip, true},
	})

	if err != nil {
		if session.UserEmail.String == "" {
			return nil, fmt.Errorf("Сессии не существует")
		}

		return nil, err
	}

	if !session.UserEmail.Valid {
		return nil, fmt.Errorf("not email")
	}

	userModel := &models.Session{
		ID:           session.ID,
		UserEmail:    session.UserEmail.String,
		RefreshToken: session.RefreshToken.String,
		AccessToken:  session.AccessToken.String,
		UserAgent:    session.UserAgent.String,
		IP:           session.Ip.String,
		ExpiresIn:    session.ExpiresIn,
		CreatedAt:    session.CreatedAt,
	}

	return userModel, nil
}

func (o *PostgresRepository) UpdateSession(
	email string, refreshToken string, newRefreshToken string, expiresIn int64) error {
	var err error

	err = o.queries.UpdateSession(context.Background(), postgres.UpdateSessionParams{
		UserEmail:      sql.NullString{email, true},
		RefreshToken:   sql.NullString{refreshToken, true},
		RefreshToken_2: sql.NullString{newRefreshToken, true},
		ExpiresIn:      expiresIn,
	})

	return err
}

func (o *PostgresRepository) UpdateSessionTokenOnly(
	refreshToken string, newRefreshToken string, expiresIn int64) error {
	var err error

	err = o.queries.UpdateSessionTokenOnly(context.Background(), postgres.UpdateSessionTokenOnlyParams{
		RefreshToken:   sql.NullString{refreshToken, true},
		RefreshToken_2: sql.NullString{newRefreshToken, true},
		ExpiresIn:      expiresIn,
	})

	return err
}

func (o *PostgresRepository) ReadEmailRoleWithRefreshToken(refreshToken string) (*models.UserEmailRole, error) {
	var err error

	emailAndRole, err := o.queries.ReadEmailRoleFromSessions(context.Background(), sql.NullString{
		String: refreshToken,
		Valid:  true,
	})

	if err != nil {
		if emailAndRole.Email == "" {
			return nil, fmt.Errorf("Пользователя не существует")
		}

		return nil, fmt.Errorf("что то пошло не так")
	}

	// Todo: в запросе уже чекать на из deleted? и возваращать время жизни
	//if emailAndRole.IsDeleted.Bool {
	//	return nil, fmt.Errorf("Пользователя не существует, или аккаунт был удален")
	//}

	userModel := &models.UserEmailRole{
		Email: emailAndRole.Email,
		Role:  emailAndRole.Role,
	}

	return userModel, err
}

func (o *PostgresRepository) DeleteSession(email string, ip, userAgent string) error {
	var err error

	err = o.queries.DeleteSession(context.Background(), postgres.DeleteSessionParams{
		UserAgent: sql.NullString{email, true},
		Ip:        sql.NullString{ip, true},
		UserEmail: sql.NullString{userAgent, true},
	})

	return err
}

func (o *PostgresRepository) CreatContent(email string, name, description string) error {
	var err error
	// Todo: нет обработки ошибок
	err = o.queries.CreateContent(context.Background(), postgres.CreateContentParams{
		UserEmail:   sql.NullString{email, true},
		Name:        sql.NullString{name, true},
		Description: sql.NullString{description, true},
	})

	return err
}

func (o *PostgresRepository) ReadContent(email string, id int32) (*models.Content, error) {
	var err error

	content, err := o.queries.ReadContent(context.Background(), postgres.ReadContentParams{
		UserEmail: sql.NullString{email, true},
		ID:        id,
	})

	// Todo: нет обработки ошибок
	if err != nil {
		return nil, err
	}

	if !content.UserEmail.Valid {
		return nil, fmt.Errorf("not email")
	}

	contentModel := &models.Content{
		ID:          content.ID,
		UserEmail:   content.UserEmail.String,
		Name:        content.Name.String,
		Description: content.Description.String,
	}

	return contentModel, err
}

func (o *PostgresRepository) CreatBlog(name, description string, img []byte) error {
	err := o.queries.CreateBlog(context.Background(), postgres.CreateBlogParams{
		Name:        sql.NullString{name, true},
		Description: sql.NullString{description, true},
		Img:         img,
	})

	val, ok := err.(*pgconn.PgError)
	if ok && pgerrcode.IsIntegrityConstraintViolation(val.Code) {
		err = fmt.Errorf("Контент уже существует")
		return err
	}

	if err != nil {
		return fmt.Errorf("что то пошло не так")
	}

	return nil
}

func (o *PostgresRepository) ReadBlog(id int32) (*models.Blog, error) {
	blog, err := o.queries.ReadBlog(context.Background(), id)
	if err != nil {
		if blog.Name.String == "" {
			return nil, fmt.Errorf("Контента не существует")
		}

		return nil, fmt.Errorf("что то пошло не так")
	}

	if !blog.Description.Valid {
		return nil, fmt.Errorf("Description is empty")
	}

	blogModel := &models.Blog{
		ID:          blog.ID,
		Name:        blog.Name.String,
		Description: blog.Description.String,
		Img:         blog.Img,
	}

	return blogModel, nil
}
