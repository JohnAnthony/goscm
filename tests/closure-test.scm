(define counter
  (begin
    (define n 0)
    (lambda ()
      (begin (set! n (+ 1 n))
             n))))

(counter)
(counter)
(counter)
(counter)
(counter)
(counter)
