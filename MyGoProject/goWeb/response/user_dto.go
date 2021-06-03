/*
创建者：     Zero
创建时间：   2021/5/26
项目名称：   goWeb
*/
package response

import "MyGoProject/model"

type UserDto struct {
	Name      string `json:"name"`
	Telephone string `json:"telephone"`
}

//DTO就是数据传输对象（Data Transfer Object）的缩写；用于展示与服务层之间的数据传输对象
func ToUserDto(user model.UserBase) UserDto {
	return UserDto{
		Name:      user.Name,
		Telephone: user.Telephone,
	}
}
