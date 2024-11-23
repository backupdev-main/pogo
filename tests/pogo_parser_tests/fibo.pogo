program recursiveFibo;

var result : int;

func fib(n : int) {
    var temp1, temp2 : int;

    if(n < 2) {
        result = n;
    } else {
        temp1 = n - 1;
        temp2 = n - 2;
        fib(temp1)
        temp1 = result;
        fib(temp2)
        temp2 = result;
        result = temp1 + temp2;
    }
};

begin
    fib(30)
    print("This is the result", result)
end

