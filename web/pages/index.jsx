import {withLayout} from "../components/layout"

const IndexPage = withLayout(() => {
    return (
        <>  
            <section>
                <h2>Getting Started</h2>
            </section>            
        </>
    )
}, {
    description: 'This page is an overview of the Orbit documentation',
    title: 'Orbit',
})

export default IndexPage