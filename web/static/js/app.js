
$(document).ready( function() {
    $(".action_drop_garden").click(function() {
        id = $(this).data("garden-id");
        $.ajax({
            url: '/api/gardens/'+id,
            type: 'DELETE',
            success: function(result) {
                // Do something with the result
                location.reload();
            }, error: function(result) {
                // Do something with the result
                alert("Failed to drop garden...");
                location.reload();
            }
        });
    });
    $(".action_water").click(function() {
        pid = $(this).data("parcel-id");
        gid = $(this).data("garden-id");
        $.ajax({
            url: '/api/gardens/'+gid+'/hydro/'+pid,
            type: 'POST',
            success: function(result) {
                // Do something with the result
                $(".hydro_status[data-parcel-id="+pid+"]").html(result.hydro)
                console.log(result)
            }, error: function(result) {
                // Do something with the result
                alert("Failed to water garden...");
                location.reload();
            }
        });
    });
    $(".action_get_water").click(function() {
        pid = $(this).data("parcel-id");
        gid = $(this).data("garden-id");
        $.ajax({
            url: '/api/gardens/'+gid+'/hydro/'+pid,
            type: 'GET',
            success: function(result) {
                // Do something with the result
                $(".hydro_status[data-parcel-id="+pid+"]").html(result.hydro)
                console.log(result)
            }, error: function(result) {
                // Do something with the result
                alert("Failed to get water garden...");
                location.reload();
            }
        });
    });
    $(".action_add_plant").click(function() {
        alert("Not implemented yet.");
    });
    
})