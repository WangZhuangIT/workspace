package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {

	ch := make(chan int)
	page := 1
	for index := 0; index < page; index++ {
		go downLoadData(index, ch)
	}

	for index := 0; index < page; index++ {
		fmt.Printf("第%v页数据爬取完毕\n", <-ch)
	}
}

func downLoadData(page int, ch chan int) {
	//https://tieba.baidu.com/f?kw=%E5%9D%A6%E5%85%8B%E4%B8%96%E7%95%8C&ie=utf-8&pn=
	// http://tieba.baidu.com/f?kw=%E5%85%B0%E5%B7%9E%E4%BA%A4%E9%80%9A%E5%A4%A7%E5%AD%A6&ie=utf-8&pn=
	doc, err := goquery.NewDocument("http://tieba.baidu.com/f?kw=%E5%85%B0%E5%B7%9E%E4%BA%A4%E9%80%9A%E5%A4%A7%E5%AD%A6&ie=utf-8&pn=" + strconv.Itoa(page*50))
	if err != nil {
		log.Fatal(err)
	}

	msg := ""
	content_data := ""
	doc.Find("ul#thread_list .j_thread_list .cleafix").Each(func(i int, contentSelection *goquery.Selection) {
		title := contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_lz .threadlist_title a").Text()
		createtime := contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_lz .threadlist_author span.is_show_create_time").Text()
		author := contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_lz .threadlist_author span.tb_icon_author").Text()
		data := strings.TrimSpace(contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_detail .threadlist_text .threadlist_abs_onlyline").Text())
		lastReplyer := strings.TrimSpace(contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_detail .threadlist_author .tb_icon_author_rely").Text())
		lastReplyerDate := strings.TrimSpace(contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_detail .threadlist_author .threadlist_reply_date").Text())
		reply := contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_left span").Text()
		link := contentSelection.Find(".j_thread_list .t_con .j_threadlist_li_right .threadlist_lz .threadlist_title a")
		href, isExist := link.Attr("href")

		if title == "" || createtime == "" || author == "" {
			return
		}

		if isExist {
			//http://tieba.baidu.com/p/5669694238
			doc, err = goquery.NewDocument("http://tieba.baidu.com" + href)
			if err != nil {
				log.Fatal(err)
			}
		}

		child_pagestr := doc.Find("div#container .content .p_thread .l_thread_info .l_posts_num .l_reply_num span[class=red]").Last().Text()
		child_page, _ := strconv.Atoi(child_pagestr)
		for index := 1; index <= child_page; index++ {
			doc, err = goquery.NewDocument("http://tieba.baidu.com" + href + "?pn=" + strconv.Itoa(index))
			if err != nil {
				log.Fatal(err)
			}

			title := strings.TrimSpace(doc.Find("div#j_p_postlist .l_post_bright .d_post_content_firstfloor .p_content_nameplate cc .j_d_post_content").Text())
			content_data += title
			//d_post_content_main
			doc.Find(".d_post_content_main").Each(func(i int, contentSelection *goquery.Selection) {
				data := contentSelection.Find(".p_content_nameplate cc .j_d_post_content").Text()
				content_data += data + "\n"
			})
		}
		fmt.Println(content_data)

		msg += "\n" + title + "\t" + author + "\t" + createtime + "\n" + data + "\n" + reply + "\t" + lastReplyer + "\t" + lastReplyerDate + "\n"
	})
	storeData(page, msg, content_data)
	ch <- page
}

func storeData(page int, msg string, content_data string) {
	fout, err := os.Create("page" + strconv.Itoa(page) + ".txt")
	if err != nil {
		fmt.Printf("打开文件失败：%v", err)
	}

	defer fout.Close()
	fout.WriteString(msg)
	fout.WriteString("\n" + content_data + "\n")
	fout.WriteString("*********************************************************************************************************************************************\n\n\n\n")
}
