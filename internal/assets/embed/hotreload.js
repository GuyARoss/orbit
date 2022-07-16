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
                    resetErrors()
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

let globalErrors = []
const resetErrors = () => {
    globalErrors = []
}

window.onerror = function(error) {
    renderOrbitErrFrame(error)
};

function renderOrbitErrFrame(error) {
    if (!!error) {
        globalErrors.push(error)
    }
    
    if (globalErrors.length === 1) {
        document.querySelector('body').innerHTML += renderErrorFrame(globalErrors)
    } else {
        document.getElementById('orbit-error-frame').innerHTML = renderErrorFrame(globalErrors, true)
    }
}

let currentErrorIdx = 0
function decrementErrorIndex() {
    if (currentErrorIdx <= 0) {
        return
    }
    currentErrorIdx -= 1
    renderOrbitErrFrame(undefined)
}

function incrementErrorIndex() {
    if (currentErrorIdx >= globalErrors.length - 1) {
        return
    }
    currentErrorIdx +=1
    renderOrbitErrFrame(undefined)
}

function closeErrorFrame() {
    document.getElementById('orbit-error-container').innerHTML = ""
}

function renderErrorFrame(errors, noParent=false) {        
    const body = `
    <div>
        <div style="display: flex;">
            <div style="margin-right: 15px;">
                <button onclick="decrementErrorIndex()">Back</button>
                <button onclick="incrementErrorIndex()">Next</button>
            </div>
            <div>
                ${currentErrorIdx + 1} of ${errors.length} unhandled errors
            </div>
        </div>

        <div style="
            margin-top: 20px;
        ">
            ${errors[currentErrorIdx]}
        </div>
    </div>
    `
    
    if(noParent) {
        return body
    }

    return `
    <div id="orbit-error-container">
        <div onclick="closeErrorFrame()" style="position: absolute; width: 100vw; height: 100vh; background: #ababab8a;">
        </div>

        <div id="orbit-error-frame" style="
            position: absolute;
            width: 500px;
            height: fit-content;
            left: 0;
            right: 0;
            margin-left: auto;
            margin-right: auto;
            padding: 20px;
            margin-top: 10%;
            box-shadow: 0px 4px 19px #d5d5d596;
            border-top: solid #e54e4e 5px;
            border-radius: 10px;
            background: white;
        ">        
            ${body}
        </div>
    </div>    
    `
}