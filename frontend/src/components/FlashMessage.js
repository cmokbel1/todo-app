export const FlashMessage = ({ messageState, returnError }) => {
    if (messageState) {
        return (
            <div className="container alert alert-success flash">
                <div className="text-center">
                    <p>{messageState}</p>
                </div>
            </div>
        )
    } else if (returnError) {
        return (
            <div className="container alert alert-danger flash">
                <div className="text-center">
                    <p>{returnError}</p>
                </div>
            </div>
        )
    }

}
