//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=../entgql/interfaces/svc/admin_user.go -destination=admin_user_mock.go -package=mocks
//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source=../entgql/interfaces/svc/role.go -destination=role_mock.go -package=mocks

package mocks
