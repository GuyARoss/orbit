import React from 'react'

import BoxArrowSvg from './box-arrow-svg'

const Layout = ({ active, children, title, description }) => {
    const activeMap = {
        index: false, 
        cli: false,
        experimental: false,
        micro: false,
        react: false,
    }

    activeMap[active] = true

    return (
        <>
            <header className="default-header">
                <h1>{title}</h1>
                <p>{description}</p>
    
                <div className="links">
                    <a href="http://github.com/GuyARoss/orbit">GitHub <BoxArrowSvg /></a> 
                    <a href="./changes.html">v0.21.0</a>
                </div>
                <hr />
            </header>
            <div className="body-content">
                <div className="sidebar">
                    <ul>
                        <li className={activeMap.index ? "active" : ""}><a href="./index.html">Getting Started</a></li>
                        <li className={activeMap.cli ? "active" : ""}><a href="./api-commands.html">CLI</a></li>         
                        <li className={activeMap.experimental ? "active" : ""}><a href="./experimental.html">Experimental Features</a></li>
                    </ul>

                    <label>Guides</label>
                    <ul>
                        <li className={activeMap.micro ? "active" : ""}><a href="./micro-frontend-guide.html">Micro-frontend</a></li>
                        <li className={activeMap.react ? "active" : ""}><a href="./react-guide.html">React</a></li>
                    </ul>  
                </div>
                <article className="main">
                    {children}
                </article>
            </div>
        </>
    )
}

export const withLayout = (Component, settings) => (props) => (
    <>
        <Layout {...settings}>
            <Component {...props}/>
        </Layout>        
    </>
)

export default Layout