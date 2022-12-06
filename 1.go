package main

//
//
//import (
//	"database/sql"
//	"fmt"
//	"github.com/antchfx/htmlquery"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"gorm.io/gorm/logger"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//// 全局数据库对象
//var db = &gorm.DB{}
//
//// 数据库配置信息
//var drive = "mysql"
//var database = "crawler"
//var host = "127.0.0.1"
//var port = 3306
//var username = "root"
//var password = "root"
//
//// 协程
//var wg sync.WaitGroup
//
//func init() {
//	// 连接数据库
//	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)
//	mysqlDb, err := sql.Open(drive, dsn)
//	if err != nil {
//		panic(err)
//	}
//	db, err = gorm.Open(mysql.New(mysql.Config{
//		Conn: mysqlDb,
//	}), &gorm.Config{
//		// 日志配置
//		Logger: logger.New(
//			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
//			logger.Config{
//				SlowThreshold:             time.Second, // 慢 SQL 阈值
//				LogLevel:                  logger.Info, // 日志级别
//				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
//				Colorful:                  false,       // 禁用彩色打印
//			}),
//	})
//	err = db.AutoMigrate(&Article{}, &ArticleCate{})
//	if err != nil {
//		return
//	}
//	if err != nil {
//		panic("数据库连接失败")
//	}
//	fmt.Println("数据库连接成功！")
//}
//
//// 文章
//type Article struct {
//	gorm.Model
//	Title    string `json:"title"`
//	CateId   int    `json:"cate_id"`
//	Status   int    `json:"status"`
//	Author   string `json:"author"`
//	SiteName string `json:"site_name"`
//	PushTime string `json:"push_time"`
//	Content  string `json:"content"`
//	Visit    int    `json:"visit"`
//	Likes    int    `json:"likes"`
//	UnLikes  int    `json:"un_likes"`
//	Level    int    `json:"level"`
//}
//
//// 文章分类
//type ArticleCate struct {
//	Id        int    `json:"id"`
//	Name      string `json:"name"`
//	Sort      int    `json:"sort"`
//	FromName  string `json:"from_name"`
//	CreatedAt time.Time
//	UpdatedAt time.Time
//}
//
//// 获取文章分类
//func getCateList(url string) (data []ArticleCate, err error) {
//	doc, err := htmlquery.LoadURL(url)
//	if err != nil {
//		return nil, err
//	}
//	list := htmlquery.Find(doc, "//ul[@id='starlist']/li")
//	var articleCate = []ArticleCate{}
//	for _, item := range list {
//		a := htmlquery.FindOne(item, "a")
//		name := htmlquery.InnerText(a)
//
//		if name != "百科网" && name != "生活知识" {
//			cate := ArticleCate{
//				Name:      name,
//				FromName:  url + htmlquery.SelectAttr(a, "href"),
//				CreatedAt: time.Now(),
//				UpdatedAt: time.Now(),
//			}
//			articleCate = append(articleCate, cate)
//		}
//	}
//	// 2、保存文章分类到数据库中
//	db.Create(&articleCate)
//	return articleCate, err
//}
//
//// 获取文章列表
//func getArticleList(baseUrl, cateUrl, siteName string, cateId int) (data []Article, err error) {
//	doc, err := htmlquery.LoadURL(cateUrl)
//	if err != nil {
//		return nil, err
//	}
//	list := htmlquery.Find(doc, "//div[@class='blogs-list']/ul/li")
//	var articleList = []Article{}
//	for _, item := range list {
//		a := htmlquery.FindOne(item, "//h2/a")
//		span := htmlquery.FindOne(item, "//span")
//		article_time := htmlquery.InnerText(span)
//
//		// 获取详情
//		detailUrl := baseUrl + htmlquery.SelectAttr(a, "href")
//		fmt.Println("详情地址：", detailUrl)
//		detail, err := getArticleDetail(detailUrl)
//		if err != nil {
//			return nil, err
//		}
//		article := Article{
//			Title:    htmlquery.InnerText(a),
//			PushTime: article_time,
//			Content:  detail,
//			CateId:   cateId,
//			SiteName: siteName,
//		}
//		articleList = append(articleList, article)
//	}
//	// 保存数据库
//
//	db.Create(&articleList)
//	return articleList, err
//}
//
//// 获取文章详情,并且保存数据库
//func getArticleDetail(url string) (string, error) {
//	doc, err := htmlquery.LoadURL(url)
//	if err != nil {
//		return "", err
//	}
//	content := htmlquery.FindOne(doc, "//div[@class='newstext']")
//	// 保存数据库
//	return htmlquery.OutputHTML(content, false), err
//}
//
//// 获取分页规则
//func getPageRule(CateUrl string) (pre string, count int, err error) {
//	doc, err := htmlquery.LoadURL(CateUrl)
//	if err != nil {
//		return "", 0, err
//	}
//	page := htmlquery.Find(doc, "//div[@class='pagelist']/a")
//	if len(page) == 0 {
//		return
//	}
//	a := htmlquery.SelectAttr(page[2], "href")
//	arr := strings.Split(a, "_")
//	// 总数
//	lastPage := htmlquery.SelectAttr(page[len(page)-1], "href")
//	arr2 := strings.Split(lastPage, "_")
//	stringTotal := strings.Split(arr2[2], ".")
//	total, _ := strconv.Atoi(stringTotal[0])
//	return fmt.Sprintf("%s_%s_", arr[0], arr[1]), total, err
//}
//
//// 开启爬取
//func spider(baseUrl, name string, cate ArticleCate, count int, pageRule string) {
//	if count > 0 {
//		for page := 1; page <= count; page++ {
//			// 获取文章列表
//			articleListUrl := fmt.Sprintf("%s%s%d.html", cate.FromName, pageRule, page)
//			fmt.Println("爬取地址：", articleListUrl)
//			list, err := getArticleList(baseUrl, articleListUrl, name, cate.Id)
//			if err != nil {
//				panic(err)
//			}
//			// 保存数据库
//			fmt.Println("列表", list)
//		}
//	}
//	wg.Done()
//}
//
//func main() {
//	baseUrl := "http://www.wenxue58.com"
//	name := "词句百科网"
//	// 1、 获取分类
//	cateData, err := getCateList(baseUrl)
//	if err != nil {
//		panic(err)
//	}
//	var articleCate = []ArticleCate{}
//	db.Where(articleCate).Order("id desc").Limit(len(cateData)).Find(&articleCate)
//	for _, cate := range articleCate {
//		// 获取分页规则
//		fmt.Println("分类：", cate)
//		pageRule, count, _ := getPageRule(cate.FromName)
//		fmt.Println(pageRule, count, name)
//		wg.Add(1)
//		go spider(baseUrl, name, cate, count, pageRule)
//	}
//	wg.Wait()
//	fmt.Println("爬取结束。。。")
//}
