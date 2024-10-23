package models

type PackSize struct {
	ID   uint32 `json:"id"`
	Size uint32 `json:"size"`
}

type PackSizeRequest struct {
	Size uint32 `json:"size" binding:"required"`
}
