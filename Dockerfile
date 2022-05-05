# не работает
FROM golang:1.16
ADD metida ./metida
ADD db/test.db ./db/test.db
EXPOSE 8080
CMD ["./metida"]