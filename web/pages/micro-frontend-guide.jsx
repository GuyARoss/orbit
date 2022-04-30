import {withLayout} from "../components/layout"

const MicrofrontendGuide = withLayout(() => {
    return (
        <>
            <header>
                <h2>Micro Frontend</h2>
                <p>You can find the code for this guide <a href="https://github.com/GuyARoss/orbit/tree/master/examples/micro-frontend">here</a></p>
            </header>
        </>
    )
}, {
    active: 'micro',
    title: 'Guide - Micro frontend',
    description: 'A simple micro frontend guide',
})


export default MicrofrontendGuide