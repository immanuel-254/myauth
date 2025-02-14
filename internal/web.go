package internal

import (
	"context"
	"net/http"
	"os"
	"strconv"

	"github.com/immanuel-254/myauth/frontend/src"
	"github.com/immanuel-254/myauth/internal/models"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	queries := models.New(DB)
	ctx := context.Background()

	todaySessions, err := queries.SessionTodayList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	yesterdaySessions, err := queries.SessionYesterdayList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	weeklySessions, err := queries.SessionWeeklyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	previousweeklySessions, err := queries.SessionPreviousWeeklyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	monthlySessions, err := queries.SessionMonthlyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	previousmonthlySessions, err := queries.SessionPreviousMonthlyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	todayLogs, err := queries.LogTodayList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	yesterdayLogs, err := queries.LogYesterdayList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	weeklyLogs, err := queries.LogWeeklyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	previousweeklyLogs, err := queries.LogPreviousWeeklyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	monthlyLogs, err := queries.LogMonthlyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	previousmonthlyLogs, err := queries.LogPreviousMonthlyList(ctx)

	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	activity := make(map[string]map[string]string)
	activity["daily"] = make(map[string]string)
	activity["weekly"] = make(map[string]string)
	activity["monthly"] = make(map[string]string)

	activity["daily"]["current"] = strconv.Itoa(len(todayLogs)) //string(len(todayLogs))
	activity["daily"]["change"] = strconv.Itoa(len(todayLogs) - len(yesterdayLogs))
	activity["weekly"]["current"] = strconv.Itoa(len(weeklyLogs))
	activity["weekly"]["change"] = strconv.Itoa(len(weeklyLogs) - len(previousweeklyLogs))
	activity["monthly"]["current"] = strconv.Itoa(len(monthlyLogs))
	activity["monthly"]["change"] = strconv.Itoa(len(monthlyLogs) - len(previousmonthlyLogs))

	Login := make(map[string]map[string]string)
	Login["daily"] = make(map[string]string)
	Login["weekly"] = make(map[string]string)
	Login["monthly"] = make(map[string]string)

	Login["daily"]["current"] = strconv.Itoa(len(todaySessions))
	Login["daily"]["change"] = strconv.Itoa(len(todaySessions) - len(yesterdaySessions))
	Login["weekly"]["current"] = strconv.Itoa(len(weeklySessions))
	Login["weekly"]["change"] = strconv.Itoa(len(weeklySessions) - len(previousweeklySessions))
	Login["monthly"]["current"] = strconv.Itoa(len(monthlySessions))
	Login["monthly"]["change"] = strconv.Itoa(len(monthlySessions) - len(previousmonthlySessions))

	component := src.Base(src.DashBoard(activity, Login))

	component.Render(context.Background(), w)
}

func Dashlogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		queries := models.New(DB)
		ctx := context.Background()

		data := map[string]string{
			"email":    r.FormValue("email"),
			"password": r.FormValue("password"),
		}

		key, code, err := AuthLogin(queries, ctx, data)

		if err != nil {
			http.Error(w, err.Error(), code)
			return
		}

		secure, err := strconv.ParseBool(os.Getenv("HTTPS"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    key,
			Path:     "/",
			HttpOnly: true,   // Prevent JavaScript access
			Secure:   secure, // Send only over HTTPS
			SameSite: http.SameSiteStrictMode,
		})

		if code == http.StatusOK {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		return
	}

	component := src.Base(src.Login())

	component.Render(context.Background(), w)
}

func Dashlogout(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		queries := models.New(DB)
		ctx := r.Context()

		var token string

		cookie, err := r.Cookie("session_token")
		if err == nil {
			token = cookie.Value // Use token from cookie if available
		}

		if token == "" {
			http.Redirect(w, r, "/dash-login", http.StatusSeeOther)
			return
		}

		session, err := queries.SessionRead(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete session
		err = queries.SessionDelete(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		Logging(queries, ctx, "session", "delete", session.ID, session.UserID, w, r)

		secure, err := strconv.ParseBool(os.Getenv("HTTPS"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // Delete immediately
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, "/dash-login", http.StatusSeeOther)
		return
	}

	component := src.Base(src.Logout())

	component.Render(context.Background(), w)
}
