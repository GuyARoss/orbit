function isHotReloadReady() {
    return document.getElementById("orbit_bk") !== null
}

function initHotReload() {
    const keys = document.getElementById("orbit_bk").attributes["src"].value.split("/")
    const primaryKey = keys[keys.length -1].replace(".js", "")
    
    // @@todo: this port number should be part of the bundle data.
    const socket = new WebSocket("ws://localhost:3005/ws");
    socket.onopen = function() {        
        socket.send(JSON.stringify({
            operation: "page",
            value: primaryKey,
        }))
    }

    socket.onmessage = function(event) {
        const incoming = JSON.parse(event.data)

        switch (incoming?.operation) {
            case "refresh": {
                window.location.reload()
            }
        }
    }
}

const interval = setInterval(() => {
    if (isHotReloadReady()) {
        clearInterval(interval)

        initHotReload()
    }
})
