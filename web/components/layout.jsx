import React from 'react'

const Layout = ({ title, description }) => {
    return (
        <>
            <header className="default-header">
                <h1>{title}</h1>
                <p>{description}</p>
    
                <div className="links">
                    <a href="http://github.com/GuyARoss/orbit">GitHub</a>
                    <a href="./changes.html">v0.7.1</a>
                </div>
                <hr />
            </header>
            
        </>
    )
}

export const withLayout = (Component, settings) => (props) => (
    <>
        <Layout {...settings}/>
        <article>
            <Component {...props}/>
        </article>
    </>
)

export default Layout