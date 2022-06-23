export const FlashMessage = ({ messageState, setMessageState, returnError }) => {
    let text;
    let classes = "flash rounded-0";
    let similarClasses = " text-center alert"

    if (messageState) {
        classes += `${similarClasses} alert-success`;
        text = messageState;
        setTimeout(() => {
            setMessageState("");
        }, 1000)
    } else if (returnError) {
        classes += `${similarClasses} alert-danger`;
        text = returnError;
    }
    return (
        <div className={classes}>
            <p>{text}</p>
        </div>
    )
}