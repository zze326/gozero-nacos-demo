// Code generated by goctl. DO NOT EDIT.
package types

type GetUserReq struct {
	ID int `path:"id"`
}

type GetUserReply struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
