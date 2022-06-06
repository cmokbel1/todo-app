

export const FlashMessage = ({ messageState }) => {
    if (messageState) {
        return (
            <div className="container mb-4" style={{ backgroundColor: 'green', color: 'snow' }}>
                <div className="text-center">
                    <p>{messageState}</p>
                </div>
            </div>
        )
    }

}
