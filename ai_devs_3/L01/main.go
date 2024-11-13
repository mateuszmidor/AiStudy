package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/mateuszmidor/AiStudy/ai_devs_3/internal/openai"
)

const loginFormURL = "http://xyz.ag3nts.org"
const loginFormExample = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>XYZ - In AI We Trust</title>
    <link rel="stylesheet" type="text/css" href="/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="/css/fontawesome-all.min.css">
    <link rel="stylesheet" type="text/css" href="/css/iofrm-style.css">
    <link rel="stylesheet" type="text/css" href="/css/iofrm-theme40.css">
</head>
<body>
    <div class="form-body without-side">
        <div class="iofrm-layout">
            <div class="form-holder">
                <div class="form-content">
                    <div class="form-items">
                        <h3 class="font-md">Xplore Your Zone</h3>
                        <p>The best patrolling robots for your factory.</p>
                        <form method="post" action="/">
                            <input class="form-control" type="text" name="username" placeholder="Login" required>
                            <input class="form-control" type="password" name="password" placeholder="Password" required>
                            <p>Prove that you are not human</p>
                            <p id="human-question">Question:<br />Rok lądowania na Księżycu?</p>
                            <input class="form-control" type="number" name="answer" placeholder="Answer" required>
                            <div class="form-button d-flex align-items-center">
                                <button id="submit" type="submit" class="btn btn-primary">Login</button><a href="forget.php">Forget password?</a>
                            </div>
                        </form>
                        <div class="page-links">
                            <a href="/register">Register new account</a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
<script src="/js/jquery.min.js"></script>
<script src="/js/popper.min.js"></script>
<script src="/js/bootstrap.bundle.min.js"></script>
<script src="/js/main.js"></script>
</body>
</html>
`

// fetchHTML retrieves the HTML content from the specified URL and returns it as a string.
func fetchHTML(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error fetching URL:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

// extractQuestionFromHTMLForm extracts the question text from a specific HTML element with id "human-question".
func extractQuestionFromHTMLForm(html string) string {
	startTag := `<p id="human-question">Question:<br />`
	endTag := `</p>`

	startIndex := strings.Index(html, startTag)
	if startIndex == -1 {
		return ""
	}
	startIndex += len(startTag)

	endIndex := strings.Index(html[startIndex:], endTag)
	if endIndex == -1 {
		return ""
	}

	return html[startIndex : startIndex+endIndex]
}

// submitLoginCredentials sends login credentials and captcha answer to the specified login page URL and returns the response body as a string.
func submitLoginCredentials(login, password, captchaAnswer, loginPageURL string) string {
	data := url.Values{}
	data.Set("username", login)
	data.Set("password", password)
	data.Set("answer", captchaAnswer)

	resp, err := http.PostForm(loginPageURL, data)
	if err != nil {
		fmt.Println("Error posting form:", err)
		return ""
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return ""
	}

	return string(body)
}

// main function orchestrates the process of fetching an HTML form, extracting a question, generating a prompt for completion, and submitting login credentials.
func main() {
	htmlForm := fetchHTML(loginFormURL)
	question := extractQuestionFromHTMLForm(htmlForm)
	prompt := "Odpowiedz maksymalnie zwięźle na ponizsze pytanie, w odpowiedzi podaj wyłącznie liczbę bez zadnego dodatkowego tekstu:\n" + question
	answer, err := openai.CompletionCheap(prompt, "", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Prompt:", prompt)
	fmt.Println("Answer:")
	fmt.Println(answer)
	response := submitLoginCredentials("tester", "574e112a", answer, loginFormURL)
	fmt.Println("")
	fmt.Println(response)
}
