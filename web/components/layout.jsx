import React from 'react'

import BoxArrowSvg from './box-arrow-svg'

const Layout = ({ children, title, description }) => {
    return (
        <>
            <header className="default-header">
                <h1>{title}</h1>
                <p>{description}</p>
    
                <div className="links">
                    <a href="http://github.com/GuyARoss/orbit">GitHub <BoxArrowSvg /></a> 
                    <a href="./changes.html">v0.7.1</a>
                </div>
                <hr />
            </header>
            <div className="body-content">
                <div className="sidebar">
                    <ul>
                        <li><a href="./index.html">Getting Started</a></li>
                        <li><a href="./micro-frontend-guide.html">Micro-frontend Guide</a></li>
                        <li><a href="./react-guide.html">React Guide</a></li>
                        <li><a href="./experimental.html">Experimental Features</a></li>
                        <li><a href="./api-commands.html">Api Commands</a></li>
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