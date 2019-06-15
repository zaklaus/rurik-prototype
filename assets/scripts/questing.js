{
    log("This is a questing demo")
    invoke("addQuest", {
        Name: "EXAMPLE"
    })

    invoke("addQuest", {
        Name: "EVENTS"
    })

    invoke("quest", {
        Name: "EVENTS",
        Args: [120]
    })

    invoke("addQuest", {
        Name: "TEST0"
    })
}
