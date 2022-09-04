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
                    resetNotices()
                    window.location.reload()
                }
                case "logger": {
                    const [logLevel, message] = incoming?.value
                    attachNoticeFrame(message, logLevel, 'compiler')
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

let globalNotices = []
const resetNotices = () => {
    globalNotices = []
}

window.onerror = function(error) {
    attachNoticeFrame(error, 2)
};

function attachNoticeFrame(message, logLevel=2, origin='client') {
    if (!!message) {
        globalNotices.push({ message, logLevel, origin })
    }
    
    if (globalNotices.length === 1) {
        document.querySelector('body').innerHTML += renderNoticeFrame({
            stack: globalNotices,
            logLevel,
            origin,
        })
    } else {
        document.getElementById('orbit-notice-frame').innerHTML = renderNoticeFrame({
            stack: globalNotices,
            noParent: true,
            logLevel,
            origin,
        })
    }
}

let currentNoticeIdx = 0
function decrementFrameIndex() {
    if (currentNoticeIdx <= 0) {
        return
    }
    currentNoticeIdx -= 1
    attachNoticeFrame(undefined)
}

function incrementFrameIndex() {
    if (currentNoticeIdx >= globalNotices.length - 1) {
        return
    }
    currentNoticeIdx +=1
    attachNoticeFrame(undefined)
}

function closeNoticeFrame() {
    document.getElementById('orbit-notice-container').innerHTML = ""
}

function renderNoticeFrame({
    stack,
    noParent=false,
}) {        
    const { message, logLevel, origin } = stack[currentNoticeIdx]

    const logStr = ({
        [0]: 'Info',
        [1]: 'Warning',
        [2]: 'Error',
    })[logLevel]

    const body = `
    <div>
        <div style="display: flex;">
            <div style="margin-right: 15px;">
                <button onclick="decrementFrameIndex()">Back</button>
                <button onclick="incrementFrameIndex()">Next</button>
            </div>
            <div>
                ${currentNoticeIdx + 1} of ${stack.length} unhandled notices
            </div>
        </div>
        <div style="
            margin-top: 10px;
            font-size: 1.3rem;
            text-transform: capitalize;
        ">
            ${logStr} - ${origin}
        </div>
        <div style="
            margin-top: 20px;
        ">
            ${message}
        </div>
    </div>
    `
    
    if(noParent) {
        return body
    }

    const color = ({
        [0]: '#336fcc',
        [1]: '#d6a828',
        [2]: '#e54e4e'
    })[logLevel]

    return `
    <div id="orbit-notice-container">
        <div onclick="closeNoticeFrame()" style="position: absolute; top: 0; left:0; width: 100vw; height: 100vh; background: #ababab8a;">
        </div>

        <div id="orbit-notice-frame" style="
            position: absolute;
            top: 0;
            width: 500px;
            height: fit-content;
            left: 0;
            right: 0;
            margin-left: auto;
            margin-right: auto;
            padding: 20px;
            margin-top: 10%;
            box-shadow: 0px 4px 19px #d5d5d596;
            border-top: solid ${color} 5px;
            border-radius: 10px;
            background: white;
        ">        
            ${body}
        </div>
    </div>    
    `
}