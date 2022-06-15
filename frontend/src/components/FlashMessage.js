export const FlashMessage = ({ messageState, setMessageState, returnError }) => {
    let text;
    let classes;

    if (messageState) {
        classes = "container alert alert-success flash";
        text = messageState;
        setTimeout(() => {
            setMessageState("");
        }, 1000)
    } else if (returnError) {
        classes = "container alert alert-danger flash";
        text = returnError;
    }
    return (
        <div className={classes}>
            <div className="text-center flash">
                <p>{text}</p>
            </div>
        </div>
    )
}