package main

func questInitEntityCommands(q *questManager) {
	/* q.registerCommand("entity", func(qs *quest, qt *questTask, args []string) bool {
		if len(args) != 2 {
			return questCommandErrorArgCount("give", qs, qt, len(args), 2)
		}

		amount, ok := qs.getNumberOrVariable(args[1])

		if !ok {
			return questCommandErrorArgType("give", qs, qt, args[1], "string", "integer")
		}

		qs.printf(qt, "giving %d of %s", amount, args[0])
		return true
	}) */
}
