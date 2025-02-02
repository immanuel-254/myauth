package internal

import (
	"context"
	"net/http"
	"strconv"

	"github.com/immanuel-254/myauth/frontend/src"
	"github.com/immanuel-254/myauth/internal/models"
)

func Dashboard(w http.ResponseWriter, r *http.Request) {
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
