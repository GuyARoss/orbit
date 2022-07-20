const Example = ({day, month, year}) => {
    return (
        <div className="orbit-integration-applied">
            <h1>Orbit!</h1>            
            <p>Welcome to this example</p>
            <p>
                Today is {day}/{month}/{year}
            </p>
        </div>        
    )
}

export default Example