import {withLayout} from "../components/layout"

const IndexPage = withLayout(() => {
    return (
        <article>  
            <header>
                <h2>Getting Started</h2>
                <p>
                    Orbit is a utility to assist with server side processing 
                    that can be used to render various web frontend tooling, 
                    requiring next to no boiler plate, using the same interface.                    
                                        
                    Additionally, Orbit can be used as a static renderer (like this website)
                    or render multiple frameworks together in a <a href="https://www.simform.com/blog/micro-frontend-architecture/"> micro frontend</a>.
                </p>

                <p>Currently orbit has support for the following tools:</p>
                <ul>
                    <li>Client-side React</li>
                    <li>Vanilla Javascript</li>
                    <li>Server-side React <span>(partial support)</span></li>
                </ul>
            </header>
            <section>
                <h2>Installation</h2>
                <p>
                    You can install orbit with the go install tool like so 
                    <span className="flag"> go install github.com/GuyARoss/orbit@latest</span>
                </p>
            </section>
            <section>
                <h2>Guides</h2>
                <p>You can find guides for the following tools listed below</p>
                <ul>
                    <li><a href="./react-guide.html">Basic React</a></li>
                    <li><a href="./micro-frontend-guide.html">Micro-frontend</a></li>
                </ul>
            </section>
        </article>
    )
}, {
    active: 'index',
    description: 'This page is an overview of the Orbit documentation',
    title: 'Orbit',
})

export default IndexPage