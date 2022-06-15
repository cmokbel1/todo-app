export const FlashMessage = ({ messageState, setMessageState, returnError }) => {
    let text;
    let classes = "flash container text-center alert ";

    if (messageState) {
        classes += "alert-success";
        text = messageState;
        setTimeout(() => {
            setMessageState("");
        }, 1000)
    } else if (returnError) {
        classes += "alert-danger";
        text = returnError;
    }
    return (
        <div className={classes}>
            <p>{text}</p>
        </div>
    )
}