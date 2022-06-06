

export const FlashMessage = ({ messageState, returnError }) => {
    if (messageState) {
        return (
            <div className="container mb-4" style={{ backgroundColor: 'green', color: 'snow' }}>
                <div className="text-center">
                    <p>{messageState}</p>
                </div>
            </div>
        )
    } else if (returnError) {
        return (
            <div className="container mb-4" style={{ backgroundColor: 'red', color: 'snow' }}>
            <div className="text-center">
                <p>{returnError}</p>
            </div>
        </div>   
        )
    }

}
