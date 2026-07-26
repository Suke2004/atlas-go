package projects

import "strings"

// filterTabClass returns CSS for a project status filter tab.
func filterTabClass(active bool) string {
	if active {
		return "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all bg-indigo-600/20 text-indigo-300 border border-indigo-500/40 shadow-sm"
	}
	return "px-4 py-2 rounded-xl text-xs font-bold uppercase tracking-wider transition-all text-slate-500 hover:text-slate-300 hover:bg-slate-800/60"
}

// viewButtonClass returns CSS for a view toggle button.
func viewButtonClass(active bool) string {
	if active {
		return "p-2 rounded-lg transition-all bg-indigo-600/20 text-indigo-300 border border-indigo-500/40"
	}
	return "p-2 rounded-lg transition-all text-slate-500 hover:text-slate-300 hover:bg-slate-800/60"
}

// statusBadgeClass returns CSS for a project status badge.
func statusBadgeClass(status string) string {
	switch strings.ToLower(status) {
	case "active":
		return "px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
	case "completed":
		return "px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase bg-indigo-500/10 text-indigo-400 border border-indigo-500/20"
	case "paused":
		return "px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase bg-amber-500/10 text-amber-400 border border-amber-500/20"
	case "archived":
		return "px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase bg-slate-500/10 text-slate-400 border border-slate-500/20"
	default:
		return "px-2.5 py-1 rounded-full text-[10px] font-mono font-bold uppercase bg-slate-500/10 text-slate-400 border border-slate-500/20"
	}
}

// milestoneTitleClass returns CSS for a milestone row — strikethrough when completed.
func milestoneTitleClass(completed bool) string {
	if completed {
		return "text-sm font-medium text-slate-500 line-through"
	}
	return "text-sm font-semibold text-white"
}

// techTagPillClass returns CSS for a tech stack tag pill.
func techTagPillClass(active bool) string {
	if active {
		return "px-2.5 py-1 rounded-lg text-[11px] font-mono font-semibold bg-indigo-600/20 text-indigo-300 border border-indigo-500/30 cursor-pointer"
	}
	return "px-2.5 py-1 rounded-lg text-[11px] font-mono font-semibold bg-slate-800/60 text-slate-400 border border-slate-700/60 hover:border-slate-600 cursor-pointer transition-colors"
}
