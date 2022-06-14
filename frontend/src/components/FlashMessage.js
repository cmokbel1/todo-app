export const FlashMessage = ({ messageState, returnError }) => {
    if (messageState) {
        return (
            <div className="container alert alert-success flash w-75">
                <div className="text-center">
                    <p>{messageState}</p>
                </div>
            </div>
        )
    } else if (returnError) {
        return (
            <div className="container alert alert-danger flash w-75">
                <div className="text-center">
                    <p>{returnError}</p>
                </div>
            </div>
        )
    }

}
