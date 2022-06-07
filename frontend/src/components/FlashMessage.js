

export const FlashMessage = ({ messageState, returnError }) => {
    if (messageState) {
        return (
            <div className="container mb-4 alert alert-success">
                <div className="text-center">
                    <p>{messageState}</p>
                </div>
            </div>
        )
    } else if (returnError) {
        return (
            <div className="container mb-4 alert alert-danger">
                <div className="text-center">
                    <p>{returnError}</p>
                </div>
            </div>
        )
    }

}
