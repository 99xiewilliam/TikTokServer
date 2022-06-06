package controller

// 测试用

//func Login(c *gin.Context) {
//	_ = c.Query("username")
//	_ = c.Query("password")
//	token, err := middleware.GenToken(1)
//	if err != nil {
//		klog.Errorf("login err:%v\n", err)
//	}
//	c.Header("token", token)
//	c.JSON(http.StatusOK, gin.H{
//		"status_code": errno.SuccessCode,
//		"status_msg":  "login success",
//		"user_id":     1,
//		"token":       token,
//	})
//	return
//}
//
//func UserInfo(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"status_code": errno.SuccessCode,
//		"status_msg":  "获取用户信息成功",
//		"user": struct {
//			Id            int64
//			Name          string
//			FollowCount   int64
//			FollowerCount int64
//			IsFollow      bool
//		}{
//			Id:            1,
//			Name:          "zhi",
//			FollowCount:   0,
//			FollowerCount: 0,
//			IsFollow:      false,
//		},
//	})
//}
