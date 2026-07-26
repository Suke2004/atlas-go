package tasks

// filterTabClass returns the CSS class for a status filter tab button.
func filterTabClass(active bool) string {
	if active {
		return "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all bg-indigo-600/20 text-indigo-300 border border-indigo-500/40 shadow-sm"
	}
	return "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all text-slate-500 hover:text-slate-300 hover:bg-slate-800/60"
}

// viewButtonClass returns the CSS class for a view toggle button (kanban / list).
func viewButtonClass(active bool) string {
	if active {
		return "p-2 rounded-lg transition-all bg-indigo-600/20 text-indigo-300 border border-indigo-500/40"
	}
	return "p-2 rounded-lg transition-all text-slate-500 hover:text-slate-300 hover:bg-slate-800/60"
}

// taskTitleClass returns the CSS class for a task title — strikes through when done.
func taskTitleClass(isDone bool) string {
	if isDone {
		return "text-sm font-medium text-slate-500 line-through"
	}
	return "text-sm font-semibold text-white leading-snug"
}

// priorityPillClass returns the CSS class for a priority badge pill.
func priorityPillClass(priority string) string {
	switch priority {
	case "critical":
		return "px-2 py-0.5 rounded-md text-[10px] font-mono font-bold uppercase bg-rose-500/10 text-rose-400 border border-rose-500/20"
	case "high":
		return "px-2 py-0.5 rounded-md text-[10px] font-mono font-bold uppercase bg-orange-500/10 text-orange-400 border border-orange-500/20"
	case "medium":
		return "px-2 py-0.5 rounded-md text-[10px] font-mono font-bold uppercase bg-amber-500/10 text-amber-400 border border-amber-500/20"
	default:
		return "px-2 py-0.5 rounded-md text-[10px] font-mono font-bold uppercase bg-slate-500/10 text-slate-400 border border-slate-500/20"
	}
}

// countTasksByStatus returns how many tasks in the list match the given status.
func countTasksByStatus(taskList []TaskWithDetails, status string) int {
	count := 0
	for _, t := range taskList {
		if t.Task.Status == status {
			count++
		}
	}
	return count
}
