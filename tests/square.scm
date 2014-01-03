(load-from-path "r5rs.scm")

(define square (lambda (x) (* x x)))
(and (= 25 (square 5))
     (= 100 (square 10))
     (= 225 (square 15)))
