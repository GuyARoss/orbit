function isHotReloadReady() {
    return document.getElementsByClassName('orbit_bk').length > 0
}

function debugData() {
    return JSON.parse(document.getElementById("debug_data").innerText)
}

async function createSocket() {    
    return new Promise((res, rej) => {
        let failCounter = 0
        const createSockInterval = setInterval(() => {
            try {
                const debug = debugData()
    
                const socket = new WebSocket(`ws://localhost:${debug?.hotReloadPort}/ws`);                
                clearInterval(createSockInterval)
    
                res(socket) 
            } catch {
                failCounter += 1
     
                if (failCounter === 5) {
                    clearInterval(createSockInterval)
    
                    rej("socket connection could not be established.")
                    return
                }
            }
        }, 500)        
    }) 
}

async function initHotReload() {
    const primaryKeys = Array.from(document.getElementsByClassName("orbit_bk")).map(x => {
        const k = x.attributes["src"].value.split("/")
        return k[k.length -1].replace(".js", "")
    })
    
    try {
        const socket = await createSocket()

        socket.onopen = function() {        
            socket.send(JSON.stringify({
                operation: "pages",
                value: primaryKeys,
            }))
        }
    
        socket.onmessage = function(event) {
            const incoming = JSON.parse(event.data)
    
            switch (incoming?.operation) {
                case "reload": {
                    window.location.reload()
                }
            }
        }
    } catch(err) {
        console.log(err)
    }    
}

const interval = setInterval(() => {
    if (isHotReloadReady()) {
        clearInterval(interval)

        initHotReload()
    }
})
