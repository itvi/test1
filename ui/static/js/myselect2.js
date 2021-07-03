$('.category').select2({
    placeholder: "选择类型",
    language: "zh-CN",
    ajax: {
        url: '/asset/categories/dropdown',
        dataType: 'json',
        data: function(p) {
            return {
                q: p.term,
                page: p.page,
                csrf_token: "{{.CSRFToken}}"
            };
        },
        processResults: function(data) {
            var r = [];
            var select2Data = $.map(data, function(o) {
                var obj = new Object();
                obj.id = o.Code;
                obj.text = o.Name;
                r.push(obj);
                return obj;
            });
            return {
                results: select2Data,
            };
        },
    }
});