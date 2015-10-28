package gcpug

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"testing"
	"time"
)

func TestSlackMessageSet(t *testing.T) {
	t.SkipNow()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	c := appengine.NewContext(req)

	api := CollectorApi{}
	items, err := api.ParseJson([]byte(`{"items":[{"tags":["google-app-engine","google-bigquery"],"owner":{"reputation":58,"user_id":4361,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/bb25fba7d94af8f3aa3e9c3bba846856?s=128&d=identicon&r=PG","display_name":"sinmetal","link":"http://ja.stackoverflow.com/users/4361/sinmetal"},"is_answered":true,"view_count":60,"accepted_answer_id":12747,"answer_count":1,"score":1,"last_activity_date":1437999924,"creation_date":1437907882,"question_id":12718,"link":"http://ja.stackoverflow.com/questions/12718/app-engine-log-table%e3%81%8b%e3%82%89view%e3%82%92%e7%94%9f%e6%88%90%e3%81%97%e3%81%9f%e3%81%84","title":"App Engine Log TableからViewを生成したい"},{"tags":["php","google-app-engine","netbeans"],"owner":{"reputation":6,"user_id":10014,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/ef307f0601bf746fa96c653ed9259c0c?s=128&d=identicon&r=PG&f=1","display_name":"kenny","link":"http://ja.stackoverflow.com/users/10014/kenny"},"is_answered":false,"view_count":50,"answer_count":1,"score":1,"last_activity_date":1436934447,"creation_date":1433522753,"last_edit_date":1433523251,"question_id":11016,"link":"http://ja.stackoverflow.com/questions/11016/netbeans8-0-2%e3%81%a7%e3%81%aegoogle-app-engine-for-php%e3%81%ae%e8%a8%ad%e5%ae%9a%e3%81%ab%e3%81%a4%e3%81%84%e3%81%a6","title":"NetBeans8.0.2でのGoogle App Engine for PHPの設定について"},{"tags":["google-app-engine","memcached"],"owner":{"reputation":1,"user_id":10032,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/b75afbba54f45121b21d7ee6785edcbe?s=128&d=identicon&r=PG&f=1","display_name":"jgt4eng5","link":"http://ja.stackoverflow.com/users/10032/jgt4eng5"},"is_answered":false,"view_count":49,"answer_count":0,"score":0,"last_activity_date":1433674813,"creation_date":1433674813,"question_id":11076,"link":"http://ja.stackoverflow.com/questions/11076/%e3%83%ad%e3%83%bc%e3%82%ab%e3%83%ab%e9%96%8b%e7%99%ba%e7%92%b0%e5%a2%83windows8-1%e3%81%8b%e3%82%89google-app-engine%e3%81%aememchache%e3%81%ab%e6%8e%a5%e7%b6%9a%e3%81%99%e3%82%8b%e6%96%b9%e6%b3%95%e3%81%af%e3%81%82%e3%82%8a%e3%81%be%e3%81%99%e3%81%8b","title":"ローカル(開発環境windows8.1)からGoogle App Engineのmemchacheに接続する方法はありますか？"},{"tags":["android","google-app-engine","oauth"],"owner":{"reputation":6,"user_id":5935,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/52f8362ecd10b2105400be8441999506?s=128&d=identicon&r=PG&f=1","display_name":"3xdjichv","link":"http://ja.stackoverflow.com/users/5935/3xdjichv"},"is_answered":false,"view_count":74,"answer_count":0,"score":1,"last_activity_date":1430787677,"creation_date":1430787677,"question_id":9828,"link":"http://ja.stackoverflow.com/questions/9828/google-%e3%81%ae%e3%82%b5%e3%83%bc%e3%83%93%e3%82%b9%e3%81%a7%e8%a8%b1%e5%8f%af%e3%81%97%e3%81%9f-oauth-%e3%82%92%e5%89%8a%e9%99%a4%e3%81%99%e3%82%8b%e6%96%b9%e6%b3%95","title":"Google のサービスで許可した OAuth を削除する方法"},{"tags":["java","mysql","google-app-engine","jdbc"],"owner":{"reputation":4,"user_id":8673,"user_type":"registered","profile_image":"http://graph.facebook.com/854583487914056/picture?type=large","display_name":"Kaname Susa","link":"http://ja.stackoverflow.com/users/8673/kaname-susa"},"is_answered":true,"view_count":133,"answer_count":2,"score":0,"last_activity_date":1425525309,"creation_date":1425446338,"question_id":7416,"link":"http://ja.stackoverflow.com/questions/7416/google-app-engineoracle-jdbc-driver-oracledriver-registermbeans-error-while-re","title":"Google App Engine:oracle.jdbc.driver.OracleDriver registerMBeans: Error while registering Oracle JDBC Diagnosability MBean"},{"tags":["php","google-app-engine"],"owner":{"reputation":48,"user_id":8013,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/2947098d845f8f7a4ee9810667551172?s=128&d=identicon&r=PG","display_name":"June Yamamoto","link":"http://ja.stackoverflow.com/users/8013/june-yamamoto"},"is_answered":true,"view_count":131,"accepted_answer_id":6944,"answer_count":1,"score":0,"last_activity_date":1424662154,"creation_date":1424658370,"question_id":6939,"link":"http://ja.stackoverflow.com/questions/6939/google-app-engine-for-php-%e3%81%a7%e3%81%ae-bigtable","title":"Google App Engine for PHP での Bigtable"},{"tags":["java","google-app-engine"],"owner":{"reputation":20,"user_id":7569,"user_type":"registered","profile_image":"https://lh5.googleusercontent.com/-U_OAcaboFJU/AAAAAAAAAAI/AAAAAAAAABw/A4A90k1WlMA/photo.jpg?sz=128","display_name":"D.Nakamura","link":"http://ja.stackoverflow.com/users/7569/d-nakamura"},"is_answered":true,"view_count":95,"accepted_answer_id":6447,"answer_count":1,"score":2,"last_activity_date":1423824878,"creation_date":1421716139,"question_id":4870,"link":"http://ja.stackoverflow.com/questions/4870/gae%ef%bc%8b%e7%8b%ac%e8%87%aa%e3%83%89%e3%83%a1%e3%82%a4%e3%83%b3%ef%bc%8bcloudgate","title":"GAE＋独自ドメイン＋CloudGate"},{"tags":["google-app-engine"],"owner":{"reputation":400,"user_id":7572,"user_type":"registered","accept_rate":50,"profile_image":"http://graph.facebook.com/1611432520/picture?type=large","display_name":"tokoi","link":"http://ja.stackoverflow.com/users/7572/tokoi"},"is_answered":true,"view_count":192,"accepted_answer_id":6429,"answer_count":1,"score":1,"last_activity_date":1423805264,"creation_date":1423799243,"question_id":6424,"link":"http://ja.stackoverflow.com/questions/6424/googleappengine-%e3%81%ab%e3%82%88%e3%82%8b%e9%9d%99%e7%9a%84%e3%83%95%e3%82%a1%e3%82%a4%e3%83%ab%e9%85%8d%e4%bf%a1%e3%81%ae%e5%be%93%e9%87%8f%e8%aa%b2%e9%87%91%e3%81%ab%e3%81%a4%e3%81%84%e3%81%a6","title":"GoogleAppEngine による静的ファイル配信の従量課金について"},{"tags":["go","google-app-engine"],"owner":{"reputation":13,"user_id":2684,"user_type":"registered","profile_image":"https://www.gravatar.com/avatar/046c6a894b9b662bc0214d0261d5bab1?s=128&d=identicon&r=PG","display_name":"keima","link":"http://ja.stackoverflow.com/users/2684/keima"},"is_answered":true,"view_count":60,"accepted_answer_id":5802,"answer_count":1,"score":2,"last_activity_date":1422897492,"creation_date":1422885770,"question_id":5794,"link":"http://ja.stackoverflow.com/questions/5794/appengine-geopoint%e3%81%abjson%e3%82%bf%e3%82%b0%e3%82%92%e4%bb%98%e4%b8%8e%e3%81%97%e3%81%9f%e3%81%84","title":"appengine.GeoPointにjsonタグを付与したい"},{"tags":["java","google-app-engine","guava"],"owner":{"reputation":20,"user_id":7569,"user_type":"registered","profile_image":"https://lh5.googleusercontent.com/-U_OAcaboFJU/AAAAAAAAAAI/AAAAAAAAABw/A4A90k1WlMA/photo.jpg?sz=128","display_name":"D.Nakamura","link":"http://ja.stackoverflow.com/users/7569/d-nakamura"},"is_answered":true,"view_count":117,"accepted_answer_id":5389,"answer_count":1,"score":1,"last_activity_date":1422325900,"creation_date":1422324350,"last_edit_date":1422325374,"question_id":5385,"link":"http://ja.stackoverflow.com/questions/5385/google-gcs-client-library-%e3%81%ae-gcsfileoptions-%e3%81%a7-java-lang-nosuchmethoderror","title":"Google GCS Client Library の GcsFileOptions で java.lang.NoSuchMethodError"},{"tags":["java","google-app-engine"],"owner":{"reputation":832,"user_id":450,"user_type":"registered","accept_rate":62,"profile_image":"http://i.stack.imgur.com/3bw8U.jpg?s=128&g=1","display_name":"Yuya Matsuo","link":"http://ja.stackoverflow.com/users/450/yuya-matsuo"},"is_answered":true,"view_count":926,"accepted_answer_id":1858,"answer_count":7,"score":1,"last_activity_date":1421821231,"creation_date":1418105919,"last_edit_date":1418110781,"question_id":478,"link":"http://ja.stackoverflow.com/questions/478/google-app-engine%e3%81%a7%e4%bd%bf%e3%81%88%e3%82%8bjava%e3%83%95%e3%83%ac%e3%83%bc%e3%83%a0%e3%83%af%e3%83%bc%e3%82%af","title":"Google App Engineで使えるJavaフレームワーク"}],"has_more":false,"quota_max":300,"quota_remaining":261}`))
	if err != nil {
		t.Fatal("fatal parse json : %s", err.Error())
	}

	s := &Stackoverflow{
		Title:            items[0].Title,
		Link:             items[0].Link,
		IsAnswered:       items[0].IsAnswered,
		ViewCount:        items[0].ViewCount,
		Score:            items[0].Score,
		Tags:             items[0].Tags,
		Owner:            items[0].Owner,
		CreationDate:     (time.Time)(items[0].CreationDate),
		LastActivityDate: (time.Time)(items[0].LastActivityDate),
	}

	sm := SlackMessage{}
	sm.Set(s)
	_, err = api.PostToSlack(c, sm)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestPostToSlack(t *testing.T) {
	t.SkipNow()

	opt := &aetest.Options{AppID: "unittest", StronglyConsistentDatastore: true}
	inst, err := aetest.NewInstance(opt)
	defer inst.Close()
	req, err := inst.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal("fatal new request error : %s", err.Error())
	}

	c := appengine.NewContext(req)

	sf := SlackField{
		Title: "たぐだよ！",
	}

	sa := SlackAttachment{
		Color:      "#36a64f",
		AuthorName: "sinmetal",
		AuthorLink: "https://twitter.com/sinmetal",
		AuthorIcon: "https://pbs.twimg.com/profile_images/552632264463376384/xJJ6FKsO.png",
		Title:      "たいとるよだ！",
		TitleLink:  "http://ja.stackoverflow.com/questions/12718/app-engine-log-table%e3%81%8b%e3%82%89view%e3%82%92%e7%94%9f%e6%88%90%e3%81%97%e3%81%9f%e3%81%84",
		Fields:     []SlackField{sf},
	}

	sm := SlackMessage{
		UserName:    "sinmetal_bot",
		IconUrl:     "https://pbs.twimg.com/profile_images/552632264463376384/xJJ6FKsO.png",
		Text:        "俺が魔王だ！",
		Attachments: []SlackAttachment{sa},
	}

	api := CollectorApi{}
	_, err = api.PostToSlack(c, sm)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
