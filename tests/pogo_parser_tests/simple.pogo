program ComprehensiveTest;

var x : int;

func factorial(num : int) {
     if (num > 5) {
            print(num)
            num = num - 1;
            factorial(num)
      }
};


begin
    factorial(7)
end