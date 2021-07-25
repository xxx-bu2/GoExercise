package main

import (
	"errors"
	p_errors "github.com/pkg/errors"
	"log"
)

/*
问：我们在数据库操作的时候，比如 dao 层中当遇到一个 sql.ErrNoRows 的时候，是否应该 Wrap 这个 error，抛给上层。为什么，应该怎么做请写出代码？
答：数据库包方法，返回golang的原生error，dao进行error的wrap操作，并一路透传至最外层，由最外层记录日志
*/

//方法入口
func main() {
	data, err := logic()
	if err != nil {
		log.Printf("logic error:%T %+v\n", p_errors.Cause(err), p_errors.Cause(err))
		log.Printf("logic error stack trace:\n%+v\n", err)
		return
	}
	log.Printf("data:%+v", data)
}

//业务逻辑层
func logic() (data interface{}, err error) {
	data, err = dao()
	if err != nil {
		return nil, err
	}
	return data, nil
}

//数据访问方法
func dao() (data interface{}, err error) {
	data, err = getRows()
	if err != nil {
		return nil, p_errors.Wrap(err, "main:dao getRows error")
	}
	return data, nil
}

//第三方包方法,数据库操作
func getRows() (data interface{}, err error) {
	return nil, errors.New("sql.ErrNoRows")
}
