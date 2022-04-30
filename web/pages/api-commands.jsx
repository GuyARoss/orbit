import {withLayout} from "../components/layout"
import Notice from '../components/notice'

const ApiCommands = withLayout(() => {
    return (
        <article>
            <header>
                <h2>Commands</h2>
                <p>Find below a list of all of orbits cli commands along with usage examples</p>
                <ul>
                    <li><a className="local" href="#build">build</a></li>
                    <li><a className="local" href="#deploy">deploy</a></li>
                    <li><a className="local" href="#dev">dev</a></li>
                    <li><a className="local" href="#experimental">experimental</a></li>
                    <li><a className="local" href="#init">init</a></li>
                    <li><a className="local" href="#tools">tool</a></li>
                </ul>
                <section id="base-flags">
                    <h3>Base Flags</h3>
                    <p>
                        The following flags are supported by all build type commands. (build, dev, init & deploy)
                    </p>
                    <ul>
                        <li><strong>--mode</strong> a string to override the provided build mode. (development/ production)</li>
                        <li><strong>--pacname</strong> a string to specify the package name for the auto generated code files </li>                    
                        <li><strong>--webdir</strong>a path that specifies the location of the input directory for orbit to bundle</li>
                        <li><strong>--out</strong> a path to write out the auto generated files to</li>
                        <li><strong>--publicdir</strong>a path for an html file, that will override the default HTML structure (note: this only overrides the head & body tags)</li>
                        <li><strong>--nodemod</strong>a path that specifies the location of node_modules</li>
                        <li><strong>--depout</strong>a path that specifies the output location of a <a className="local" href="#tool-dependgraph">dependency map</a></li>
                        <li><strong>--experimental</strong>command delimited string specifying a list of experimental features <a href="./experimental.html">List of experimental features</a></li>                
                    </ul>
                </section>
            </header>            
            
            <section id="build">
                <h2>Build</h2>
                <hr />
                <p>The build command bundles the provided pages in production mode.</p>
                <h3>Flags</h3>
                <p>In addition to the base flags, the build command supports the following:</p>
                <ul>
                    <li><strong>--auditpage</strong>a path that specifies the output of an audit file</li>
                </ul>
            </section>
            
            <section id="deploy">
                <h2>Deploy</h2>
                <hr />                
                <p>The deploy command prepares the provided pages for deployment</p>                
                <Notice
                    title="Notice"
                    description="This deployment command currently only supports static file generation."
                />
                <h3>Flags</h3>
                <p>In addition to the <a className="local" href="#base-flags">base flags</a>, the deploy command supports the following:</p>
                <ul>
                    <li><strong>--staticout</strong>a path that specifies the directory to output the static files to.</li>
                </ul>
            </section>
            <section id="dev">
                <h2>Dev</h2>
                <hr />
                <p>The dev command provides the ability to hot reload development code</p>
                <h3>Flags</h3>
                <p>In addition to the <a className="local" href="#base-flags">base flags</a>, the dev command supports the following:</p>
                <ul>
                    <li><strong>--timeout</strong>duration in milliseconds until a change will be detected <span>default: 2000</span></li>
                    <li><strong>--samefiletimeout</strong>specifies the timeout duration in milliseconds until a change will be detected for repeating files <span>default: 2000</span></li>
                    <li><strong>--hotreloadport</strong>port used for hotreloading <span>default: 3005</span></li>
                </ul>
            </section>
            <section id="experimental">
                <h2>Experimental</h2>
                <hr />
                <p>Displays a list of experimental features to be used where <span className="flag">--experimental</span> flag is accepted.</p>
                <p>You can find a copy of that list <a href="./experimental.html">here</a></p>
            </section>
            <section id="init">
                <h2>Init</h2>
                <hr />
                <p>Creates a template project and installs the required dependencies</p>
                <h3>Flags</h3>
                <p>In addition to the <a className="local" href="#base-flags">base flags</a>, the init command supports no other flags.</p>
            </section>
            <section id="tools">
                <h2>Tools</h2>
                <hr />
                <p>
                    Below is a list of tools included within orbit, you may use them like so: <span class="flag">orbit tool 'tool_name'</span>
                </p>
                <ul>
                    <li><a className="local" href="#tool-dependgraph">dependgraph</a> tool for dependency graph visualization</li>
                </ul>
                <h3>Flags</h3>
                <p>No flags are currently supported by the flags command</p>
            </section>
            <section id="tool-dependgraph">
                <h2>DependGraph</h2>
                <p>Visualization tool to display the output of the "--depout" flag.</p>
                <h3>Flags</h3>
                <ul>
                    <li><strong>--graph</strong>visualization method (avsd or dracula)</li>
                </ul>
            </section>
        </article>
    )
}, {
    active: 'cli',
    title: 'CLI',
    description: 'Overview of the the CLI commands within Orbit',
})


export default ApiCommands