(ns user)



(def pcp-summary
  {;; patient id 1
   1 {:patient-id 1
      :encounters [{:encounter-id 101 :time "2024-01-01T10:00:00Z"}
                   {:encounter-id 102 :time "2024-02-01T10:00:00Z"}]
      :last-encounter-time "2024-02-01T10:00:00Z"}
   ;; patient id 2
   2 {:patient-id 2
      :encounters [{:encounter-id 201 :time "2024-03-01T10:00:00Z"}]
      :last-encounter-time "2024-03-01T10:00:00Z"}
   ;; patient id 3
   3 {:patient-id 3
      :encounters [{:encounter-id 301 :time "2024-01-15T10:00:00Z"}
                   {:encounter-id 302 :time "2024-02-15T10:00:00Z"}]
      :last-encounter-time "2024-03-15T10:00:00Z"}})

(->> pcp-summary
     vals ;; get the vals of the map, ignore keys
     (sort-by
      (juxt
       #(count (get % :encounters)) ;; asc sort by number of encounters
       :last-encounter-time)) ;; secondary asc sort by last-encounter-time
     last ;; take the last item, corresponding to largest # of docs and latest encounter
     )