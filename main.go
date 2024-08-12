package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)

	for i := 1; i <= 5; i++ {
		//TIP You can try debugging your code. We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>. To start your debugging session,
		// right-click your code in the editor and select the <b>Debug</b> option.
		fmt.Println("i =", 100/i)
	}

	const clientID = "<your client id>"
	const clientSecret = "<your client secret>"

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// We will be using `httpClient` to make external HTTP requests later in our code
	httpClient := http.Client{}

	// Create a new redirect route route
	http.HandleFunc("/oauth/redirect", func(w http.ResponseWriter, r *http.Request) {
		// First, we need to get the value of the `code` query param
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		// Next, lets for the HTTP request to call the github oauth enpoint
		// to get our access token
		reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
		req, err := http.NewRequest(http.MethodPost, reqURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("accept", "application/json")

		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()

		// Parse the request body into the `OAuthAccessResponse` struct
		var t OAuthAccessResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Finally, send a response to redirect the user to the "welcome" page
		// with the access token
		w.Header().Set("Location", "/welcome.html?access_token="+t.AccessToken)
		w.WriteHeader(http.StatusFound)
	})

	http.ListenAndServe(":8080", nil)
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.

type OAuthAccessResponse struct {
	AccessToken string `json:"access_token"`
}
