import {withLayout} from "../components/layout"

const Experimental = withLayout(() => {
    return (
        <article>
            <header>
                <h2>Supported Experimental Features</h2>
                <p>The following experimental features are supported</p>
                <ul>
                    <li><strong>ssr</strong> uses the ssr version of a specified framework if supported web frontend tool is detected.</li>
                </ul>
            </header>
        </article>
    )
}, {
    active: 'experimental',
    title: 'Experimental Features',
    description: 'A list of experimental features in orbit',
})


export default Experimental