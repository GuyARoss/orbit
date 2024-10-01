const Example = () => {
    const date = new Date();

    return (
        <div className="orbit-integration-applied">
            <h1>Orbit!</h1>            
            <p>Welcome to this example</p>
            <p>
                Today is {date.toString()}
            </p>
        </div>        
    )
}

export default Example