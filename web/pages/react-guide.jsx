import {withLayout} from "../components/layout"

const ReactGuide = withLayout(() => {
    return (
        <>
            <header>
                <h2>React</h2>
                <p>You can find the code for this guide <a href="https://github.com/GuyARoss/orbit/tree/master/examples/basic-react">here</a></p>
            </header>
        </>
    )
}, {
    active: 'react',
    title: 'Guide - React',
    description: 'A simple react guide',
})


export default ReactGuide