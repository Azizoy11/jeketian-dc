package model

type LiveShowroomResponses struct {
	OnLives []struct {
		Lives []struct {
			RoomUrlKey       string `json:"room_url_key,omitempty"`
			StartedAt        int    `json:"started_at,omitempty"`
			RoomId           int    `json:"room_id,omitempty"`
			Image            string `json:"image,omitempty"`
			ViewNum          int    `json:"view_num,omitempty"`
			MainName         string `json:"main_name,omitempty"`
			PremiumRoomType  int    `json:"premium_room_type,omitempty"`
			StreamingUrlList []struct {
				Url string `json:"url"`
			} `json:"streaming_url_list"`
		} `json:"lives"`
	} `json:"onlives"`
}

type LiveShowroomStreamingUrlResponses struct {
	StreamingUrlList []struct {
		Url string `json:"url,omitempty"`
	} `json:"streaming_url_list,omitempty"`
}
