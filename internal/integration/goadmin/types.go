package goadmin

// Envelope describes the stable JSON envelope returned by go-admin.
type Envelope[T any] struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Result  T      `json:"result"`
}

// LoginResult is the token pair returned by go-admin login endpoints.
type LoginResult struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

// LoginUserInfo is the login user payload returned by go-admin.
type LoginUserInfo struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Avatar    string `json:"avatar"`
	Sex       string `json:"sex"`
	Email     string `json:"email"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// LoginEnvelope is the full response returned by go-admin login endpoint.
type LoginEnvelope struct {
	Code     int           `json:"code"`
	Msg      string        `json:"msg"`
	Message  string        `json:"message"`
	Result   LoginResult   `json:"result"`
	UserInfo LoginUserInfo `json:"userInfo"`
}

// ProfileResult is the current-user payload returned by go-admin.
type ProfileResult struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Sex       string `json:"sex"`
	Phone     string `json:"phone"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// MenuResult is one menu item returned by go-admin.
type MenuResult struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ParentID      int64  `json:"parent_id"`
	Path          string `json:"path"`
	WebIcon       string `json:"web_icon"`
	Sort          int    `json:"sort"`
	Level         int    `json:"level"`
	ComponentName string `json:"component_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// CurrentMenusResult wraps the current-user menus payload returned by go-admin.
type CurrentMenusResult struct {
	List []MenuResult `json:"list"`
}
