(define fact
  (lambda (n)
    (if (= n 1)
        n
        (* n (fact (- n 1))))))

(fact 25)
