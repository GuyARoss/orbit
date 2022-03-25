import Thing2 from '../components/thing2'

const Example = ({day, month, year}) => {
    return (
        <>
            <h1>Orbit-SSR</h1>
            <p>Welcome to this example!</p>
                <Thing2 />
            <p>
                Today is {day}/{month}/{year}
            </p>
        </>
    )
}

export default Example