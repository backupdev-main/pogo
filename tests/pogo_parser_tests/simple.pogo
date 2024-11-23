program Factorial;

var result : float;
var x : int;

begin
    // Comments work as well!!!!
    /*
        multiline comments also work!
        look!
    */
    result = 1;
    x = 5;
    while (x > 0) {
        result = result * x;
        x = x - 1;
    }
    print("This is the result", result)
end