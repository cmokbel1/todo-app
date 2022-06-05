

export const FlashMessage = ({ errorMessage, successMessage }) => {

    if (errorMessage) {
        return (
            <div className="container mb-4" style={{ backgroundColor: 'red', color: 'snow' }}>
                <div className="text-center">
                    <p>{errorMessage}</p>
                </div>
            </div>
        )
    } else if (successMessage) {
        return (
            <div className="container mb-4" style={{ backgroundColor: 'green', color: 'snow' }}>
                <div className="text-center">
                    <p>{successMessage}</p>
                </div>
            </div>
        )
    }
}
