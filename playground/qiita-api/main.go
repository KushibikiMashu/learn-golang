package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// .env を読み込む
	err := godotenv.Load("./.env")
	if err != nil {
		panic(err)
	}

	authKey := os.Getenv("QIITA_AUTH_KEY")

	// request 作成
	url := "https://qiita.com/api/v2/authenticated_user"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+authKey)

	// client 作成
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	type User struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	var user User

	if err := json.Unmarshal(body, &user); err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))
	fmt.Println("----")
	fmt.Println(user.Id)
	fmt.Println(user.Name)
	fmt.Println(user.Description)
}

// 結果
// {
//     "description":"技術ブログに移行しました  https://panda-program.com/\r\n\r\nOOPとTDDとペア・ct+TypeScriptも好き。 ",
//     "facebook_id":"",
//     "followees_count":11,
//     "followers_count":24,
//     "github_l_name":null,
//     "id":"Panda_Program",
//     "items_count":0,
//     "linkedin_id":"",
//     "location":"Tokyo",
//     "name":"プログラミングをするパンダ",
//     "organization":"弁護士ドットコム",
//     "permanent_id":229830,
//     "profil/qiita-image-store.s3.ap-northeast-1.amazonaws.com/0/229830/profile-images/1557580795",
//     "team_only":false,
//     "twitter_screen_name":"Panda_Program",
//     "website_url":"https://panda-program.com/",
//     "image_monthly_upload_limit":104857600,
//     "image_monthly_upload_remaining":104857600
// }
// ----
// Panda_Program
// プログラミングをするパンダ
// 技術ブログに移行しました  https://panda-program.com/
//
// OOPとTDDとペア・モブプロが好きなエンジニア。React+TypeScriptも好き。
